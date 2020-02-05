
$arch = "linux", "windows"
# https://github.com/golang/go/wiki/WindowsCrossCompiling
# GOOS=windows GOARCH=386 go build -o hello.exe hello.go
# 
foreach ($envVar in $arch) {
    
    $env:GOOS = $envVar
    if ($envVar -eq "windows") {
        go build -o "hausarbeit_eb_$($envVar).exe"
    }
    else {
        go build -o "hausarbeit_eb_$($envVar)"
    }
}