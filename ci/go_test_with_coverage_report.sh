#!/usr/bin/env bash

# The MIT License (MIT)

# Copyright (c) 2015-2019 GitLab B.V.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

# from: https://gitlab.com/gitlab-org/gitlab-runner

set -eo pipefail

testsDefinitions="testsdefinitions.txt"

TESTFLAGS=${TESTFLAGS:-"-cover"}
PARALLEL_TESTS_LIMIT=${PARALLEL_TESTS_LIMIT:-1}

CI_NODE_TOTAL=${CI_NODE_TOTAL:-1}
CI_NODE_INDEX=${CI_NODE_INDEX:-0}

output="regular"
coverMode="count"
coverPkg="./..."
srcPkg="github.com/atlarge-research/opendc-emulate-kubernetes"

if [[ ${TESTFLAGS} = *"-race"* ]]; then
    output="race"
    coverMode="atomic"
fi

printMessage() {
  echo -e "\\033[1m${*}\\033[0m"
}

joinBy() {
    local IFS="${1}"
    shift
    echo "${*}"
}

prepareTestCommands() {

    local definitions

    [[ ! -f ${testsDefinitions} ]] || rm ${testsDefinitions}

    for pkg in ${OUR_PACKAGES}; do
        local testIndex=0
        local runTests=()

        echo "Listing tests for ${pkg} package"

        local tempFile
        tempFile=$(mktemp)
        local exitCode=0

        go test -list "Test.*" "${pkg}" > "${tempFile}" || exitCode=1

        local tests
        tests=$(grep "^Test" "${tempFile}" || true)

        rm "${tempFile}"

        if [[ ${exitCode} -ne 0 ]]; then
            exit ${exitCode}
        fi

        if [[ -z "${tests}" ]]; then
            continue
        fi

        local counter=0
        for test in ${tests}; do
            counter=$((counter+1))
            runTests+=("${test}")

            if [[ ${counter} -ge ${PARALLEL_TESTS_LIMIT} ]]; then
                if [[ ${#runTests[@]} -gt 0 ]]; then
                    definitions=$(joinBy "|" "${runTests[@]}")
                    echo "${pkg} ${testIndex} ${definitions}" | tee -a ${testsDefinitions}
                fi

                counter=0
                runTests=()

                testIndex=$((testIndex+1))
            fi
        done

        if [[ ${#runTests[@]} -gt 0 ]]; then
            definitions=$(joinBy "|" "${runTests[@]}")
            echo "${pkg} ${testIndex} ${definitions}" | tee -a ${testsDefinitions}
        fi

    done

	# Randomize order
	shuf ${testsDefinitions} > ${testsDefinitions}.tmp
    mv ${testsDefinitions}.tmp ${testsDefinitions}

	readarray -t e2e <<<$(cat ${testsDefinitions} | grep e2e)
	readarray -t non_e2e <<<$(cat ${testsDefinitions} | grep -v e2e)
	
	non_e2e_length=${#non_e2e[@]}
	e2e_length=${#e2e[@]}
	div=$((non_e2e_length/e2e_length))
	res=${div%.*}

	local out=()

	for i in ${!non_e2e[@]}; do

		out+=("${non_e2e[$i]}")
		if [[ $((i % res)) -eq 0 ]]; then 
			out+=("${e2e[$(($i/$res))]}")
		fi
	done

	rm ${testsDefinitions}
	for i in ${!out[@]}; do
		echo ${out[i]} >> ${testsDefinitions}
	done

    sed -i '/^$/d' ${testsDefinitions} 
}

executeTestCommand() {
    local pkg=${1}
    local index=${2}
    local runTestsList=${3}

    local options=""

    local pkgSlug
    pkgSlug=$(echo "${pkg}" | tr "/" "-")

    if [[ ${TESTFLAGS} = *"-cover"* ]]; then
        mkdir -p ".cover"
        mkdir -p ".testoutput"

        printMessage "\\n\\n--- Starting part ${index} of go tests of '${pkg}' package with coverprofile in '${coverMode}' mode:\\n"

        local profileFile=".cover/${pkgSlug}.${index}.${coverMode}.cover.txt"
        options="-covermode=${coverMode} -coverprofile=${profileFile} -coverpkg=${coverPkg}"
    else
        echo "Starting go test"
    fi

    local testOutputFile=".testoutput/${pkgSlug}.${index}.${output}.output.txt"

    local exitCode=0
    # shellcheck disable=SC2086
    go test ${options} ${TESTFLAGS} -v "${pkg}" -run "${runTestsList}" 2>&1 | tee "${testOutputFile}" || exitCode=1

    return ${exitCode}
}

executeTestPart() {
    rm -rf ".cover/"
    rm -rf ".testoutput/"

    local numberOfDefinitions
    numberOfDefinitions=$(< "${testsDefinitions}" wc -l)
    local executionSize
    executionSize=$((numberOfDefinitions/CI_NODE_TOTAL+1))
    local nodeIndex=$((CI_NODE_INDEX-1))
    local executionOffset
    executionOffset=$((nodeIndex*executionSize+1))

    printMessage "Number of definitions: ${numberOfDefinitions}"
    printMessage "Suite size: ${CI_NODE_TOTAL}"
    printMessage "Suite index: ${CI_NODE_INDEX}"

    printMessage "Execution size: ${executionSize}"
    printMessage "Execution offset: ${executionOffset}"

    local exitCode=0
    while read -r pkg index tests; do
        executeTestCommand "${pkg}" "${index}" "${tests}" || exitCode=1
    done < <(tail -n +${executionOffset} ${testsDefinitions} | head -n ${executionSize})

    exit ${exitCode}
}

computeCoverageReport() {
    local reportDirectory="out/coverage"
    local sourceFile="${reportDirectory}/coverprofile.${output}.source.txt"
    local htmlReportFile="${reportDirectory}/coverprofile.${output}.html"
    local textReportFile="${reportDirectory}/coverprofile.${output}.txt"

    mkdir -p "${reportDirectory}"

    echo "mode: ${coverMode}" > ${sourceFile}
    grep -h -v -E -e "^mode:" -e "\/mock_[^\.]+\.go" .cover/*.${coverMode}.cover.txt | grep -e "${srcPkg}" >> ${sourceFile}

	# Exclude generated code from coverage report
	sed -i '/mock/d' ${sourceFile}
	sed -i '/pb/d' ${sourceFile}
	sed -i '/deepcopy/d' ${sourceFile}

	# Exclude docker code from CI coverage, it is expected to be run locally
	sed -i '/docker/d' ${sourceFile}

    printMessage "Generating HTML coverage report"
    go tool cover -o ${htmlReportFile} -html=${sourceFile}
    printMessage "Generating TXT coverage report"
    go tool cover -o ${textReportFile} -func=${sourceFile}

    printMessage "General coverage percentage:"
    local total
    total=$(grep "total" "${textReportFile}" || echo "")

    if [[ -n "${total}" ]]; then
        echo "${output} ${total}"
    fi
}

computeJUnitReport() {
    local reportDirectory="out/junit"
    local concatenatedOutputFile="/tmp/test-output.txt"

    mkdir -p "${reportDirectory}"

    cat .testoutput/*.${output}.output.txt > "${concatenatedOutputFile}"

    go-junit-report < "${concatenatedOutputFile}" > "${reportDirectory}/report.xml"
}

case "$1" in
    prepare)
        prepareTestCommands
        ;;
    execute)
        executeTestPart
        ;;
    coverage)
        computeCoverageReport
        ;;
    junit)
        computeJUnitReport
        ;;
esac
