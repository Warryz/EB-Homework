# This will be written to the file name.
$testCase = "three_tier"

$numberOfUsersToSimulate = 10, 50, 100, 200, 400, 800, 1500
$maxNumerOfRequests = 50000

# AWS Hostname
$ipAddress = Read-Host -Prompt "Please enter the AWS Hostname"

foreach($users in $numberOfUsersToSimulate){
    $testWithUsers = $testCase + "_" + $users
    <#
    Parameters: 
    C: Number of users to simulate.
    R: Number of created users per second.
    N: Maximum number of requests to simulate(End).
    #>
    
    # Executing the Python file.
    locust -f .\Load_Test.py --no-web -c $users -r 10 -n $maxNumerOfRequests --csv=..\Locust-Ergebnisse\$testWithUsers --host="http://$($ipAddress):10000"
}
