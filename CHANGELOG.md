# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

##  [Unreleased]
N/A


##  [0.2.0]
###  Added
- Added a simple "slot filling" mechanism using the context.
- Added field ``context`` to the Response and Request objects.
- Added field ``auto_reprompt`` to the Response object.

###  Changed
- NLU object is now defined inside the SDK. (no oratio dependencies)
- NLG object is now defined inside the SDK. (no oratio dependencies)


##  [0.1.0]
###  Added
- Init repository with request and response object
- The server object handle the marshalling and unmarshalling of objects