# stopgo
Static site / pdf generator for your resume/cv, written in go

stopgo = resume.  Get it?

This code will take a resume definition in JSON format and generate a PDF and a static website.

# Example

Generated website:  http://johntunison.info
Generated PDF:  http://johntunison.info/John_Tunison_Resume.pdf

# Features
* QR code in PDF points to website
* Publish to S3 (with diffing md5 hashes, so it does quick updates)
