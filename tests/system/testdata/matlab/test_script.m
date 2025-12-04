% Copyright 2025 The MathWorks, Inc.

% Test script for MATLAB MCP Server
% This script exercises various MATLAB capabilities through the MCP server
% Used by: run_matlab_file tool test
%
% SELF-VALIDATING PATTERN: This file runs all tests and reports final status.
% System tests check for "test_script.m: ALL TESTS PASSED" at the end.

fprintf('\n=== Starting test_script.m ===\n\n');

failed_tests = {};

% Test 1: Array operations and vectorization
try
    x = 1:10;
    y = x.^2;
    assert(length(y) == 10, 'Array length mismatch');
    fprintf('  PASS: Array operations\n');
catch e
    fprintf('  FAIL: Array operations: %s\n', e.message);
    failed_tests{end+1} = 'Array operations';
end

% Test 2: Statistical functions
try
    mean_val = mean(y);
    max_val = max(y);
    min_val = min(y);
    
    assert(mean_val > 0, 'Mean should be positive');
    assert(max_val >= mean_val, 'Max should be >= mean');
    assert(min_val >= 0, 'Min should be non-negative');
    assert(max_val > min_val, 'Max should be > min for non-constant data');
    
    fprintf('  PASS: Statistics calculation\n');
catch e
    fprintf('  FAIL: Statistics calculation: %s\n', e.message);
    failed_tests{end+1} = 'Statistics calculation';
end

% Test 3: Graphics (headless)
try
    figure('Visible', 'off');
    plot(x, y, 'b-o', 'LineWidth', 2);
    xlabel('x values');
    ylabel('y = x^2');
    title('Test Plot: Quadratic Function');
    grid on;
    close(gcf);
    
    fprintf('  PASS: Figure creation\n');
catch e
    fprintf('  FAIL: Figure creation: %s\n', e.message);
    failed_tests{end+1} = 'Figure creation';
end

% Test 4: String operations
try
    result_str = sprintf('Processed %d data points', length(x));
    assert(~isempty(result_str), 'String formatting failed');
    assert(contains(result_str, '10'), 'String should contain data count');
    
    fprintf('  PASS: String operations\n');
catch e
    fprintf('  FAIL: String operations: %s\n', e.message);
    failed_tests{end+1} = 'String operations';
end

% Final summary
fprintf('\n');
if isempty(failed_tests)
    fprintf('test_script.m: ALL TESTS PASSED\n');
else
    fprintf('test_script.m: SOME TESTS FAILED (%d failures)\n', length(failed_tests));
    for i = 1:length(failed_tests)
        fprintf('  - %s\n', failed_tests{i});
    end
end
