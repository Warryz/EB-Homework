param(
    $Hostname
)

# This will be written to the file name.
$testCase = "three_tier"

$numberOfUsersToSimulate = 50, 100, 200, 400, 800, 1500

# AWS Hostname
if ($null -eq $Host) {
    $Hostname = Read-Host -Prompt "Please enter the AWS Hostname"
}

foreach ($users in $numberOfUsersToSimulate) {
    $testWithUsers = $testCase + "_" + $users
    <#
    Parameters: 
    C: Number of users to simulate.
    R: Number of created users per second.
    N: Maximum number of requests to simulate(End).
    #>
    
    # Executing the Python file.
    locust -f .\Locust\Load_Test.py --no-web -c $users -r 10 --step-load --step-clients ($users/10) --step-time 15s -t 3m --csv=Ergebnisse/$testWithUsers --host="http://$($Hostname):10000" --only-summary
}
