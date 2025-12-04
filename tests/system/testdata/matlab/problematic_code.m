function result = problematic_code(x)
    % Copyright 2025 The MathWorks, Inc.
    
    % This function has various code quality issues that checkcode should detect
    % Used to validate the check_matlab_code tool finds common problems
    
    % Issue 1: Unused variable
    unused_var = 10;
    
    % Issue 2: Variable used before defined in some code paths
    if x > 0
        y = x + 1;
    end
    result = y * 2;  % y might not be defined if x <= 0
    
    % Issue 3: Missing semicolon (intentional - causes output clutter)
    z = x + 5
    
    % Issue 4: Using eval (security/performance issue)
    eval('temp = x + 1;');
    
    % Issue 5: Inefficient loop that could be vectorized
    sum_val = 0;
    for i = 1:length(x)
        sum_val = sum_val + x(i);
    end
    
    % Issue 6: Growing array in loop (performance issue)
    data = [];
    for i = 1:10
        data(i) = i^2;
    end
end
