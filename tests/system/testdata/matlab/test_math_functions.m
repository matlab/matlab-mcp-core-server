function tests = test_math_functions
    % Copyright 2025 The MathWorks, Inc.
    
    % Test suite for basic math functions
    % Used by: run_matlab_test_file tool test
    % Validates that MATLAB test framework integration works correctly
    tests = functiontests(localfunctions);
end

function testAddition(testCase)
    % Test basic addition
    result = 2 + 2;
    expected = 4;
    verifyEqual(testCase, result, expected);
end

function testMultiplication(testCase)
    % Test basic multiplication
    result = 3 * 4;
    expected = 12;
    verifyEqual(testCase, result, expected);
end

function testSquareRoot(testCase)
    % Test square root function
    result = sqrt(16);
    expected = 4;
    verifyEqual(testCase, result, expected);
end

function testMatrixMultiplication(testCase)
    % Test matrix multiplication (validates array operations)
    A = [1 2; 3 4];
    B = [2 0; 1 2];
    result = A * B;
    expected = [4 4; 10 8];
    verifyEqual(testCase, result, expected);
end

function testMean(testCase)
    % Test statistical function
    data = [1, 2, 3, 4, 5];
    result = mean(data);
    expected = 3;
    verifyEqual(testCase, result, expected);
end

function testVectorizedOperations(testCase)
    % Test vectorized operations (common in MATLAB)
    x = 1:5;
    result = x.^2;
    expected = [1, 4, 9, 16, 25];
    verifyEqual(testCase, result, expected);
end

function testLogicalIndexing(testCase)
    % Test logical indexing (MATLAB-specific feature)
    data = [1, -2, 3, -4, 5];
    positive = data(data > 0);
    expected = [1, 3, 5];
    verifyEqual(testCase, positive, expected);
end
