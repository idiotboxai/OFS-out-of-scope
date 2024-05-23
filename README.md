# OFS - OUT OF SCOPE TOOL
Tried manually removing recusive subdomins one by one not anymore outofscope is a tool used for removing the subdomains which are out of the scope. <br> &nbsp; <br>
<img src="https://github.com/idiotboxai/OFS-out-of-scope/blob/main/logo.png" alt="logo">
 <br> &nbsp; <br>
## Installation 

```
 go install github.com/idiotboxai/OFS-out-of-scope@latest
```


## Usage
```
Usage: ofc -i <subdomains_file> -s <out_of_scope_file> -o <output_file>
```

<strong>Demo for outscope.txt</strong>
```
*.target.com
target.com/*
target.com 
```
*.target.com remove all subdomaoins and recursive subdomains to it 
if only sub is blocked target.com


## Example<br>
**subdomain.txt**
```
ab.com
abc.com
aba.com
lol.aba.com
```
**outofscope file**<br>
out.txt
```
*.aba.com
```
**outputfile**<br>
output.txt
```
ab.com
abc.com
```
