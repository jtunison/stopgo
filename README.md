# stopgo
Static site / pdf generator for your resume/cv, written in go

stopgo = resume.  Get it?

This code will take a resume definition in JSON format and generate a PDF and a static website.

# Example

* Website:  http://johntunison.info
* PDF:  http://johntunison.info/John_Tunison_Resume.pdf

# Features
* QR code in PDF points to website
* Publish to S3 (with diffing md5 hashes, so it does quick updates)

# License & Credits

stopgo is released under the Apache License. If you find it useful, please keep a link back to this page! Patches welcome, of course.

Credits:
* gofpdf - https://code.google.com/p/gofpdf/ (MIT License)
* HTML5 CSS3 site template by [Xiaoying Riley](https://www.linkedin.com/in/xiaoying) at [3rd Wave Media](http://themes.3rdwavemedia.com/); licensed under [CC BY 3.0](http://creativecommons.org/licenses/by/3.0/).
