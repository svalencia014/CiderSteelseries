Write-Output "Building Cider SteelSeries"
Start-Sleep -s 3
#Clear directory
Remove-Item ./build/* -Recurse -Force


Write-Output "Bundling using Webpack"
webpack
Write-Output "Bundling successful!"
Start-Sleep -s 3

Write-Output "Generating Blob"
node --experimental-sea-config CiderSS-config.json
Write-Output "Blob generated!"
Start-Sleep -s 3

Write-Output "Creating copy of node executable and renaming it"
Copy-Item (Get-Command node).Source ./build/CiderSS.exe
Write-Output "Copy created!"
Start-Sleep -s 3

Write-Output "Remove Signature"
signtool remove /s ./build/CiderSS.exe
Write-Output "Signature removed!"
Start-Sleep -s 3

Write-Output "Injecting cider.js into executable"
npx postject ./build/CiderSS.exe NODE_SEA_BLOB ./build/CiderSS-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2
Write-Output "Injection complete!"
Start-Sleep -s 3

Remove-Item ./build/CiderSS-prep.blob
Remove-Item ./build/CiderSS.js
Write-Output "Build Complete!"