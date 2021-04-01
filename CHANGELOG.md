# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> Deprecated
> For versions greater than 0.3.0 check github releases for release notes.

##  [0.3.0]
###  Added
- Prometheus endpoint served on ``/metrics``.

##  [0.2.1]
###  Added
- Add ``device.state`` information in the request to take some dynamic informations from device.

##  [0.2.0]
###  Added
- Added a simple "slot filling" mechanism using the context.
- Added field ``context`` to the Response and Request objects.
- Added field ``auto_reprompt`` to the Response object.

###  Changed
- Use logrus as logging system.
- NLU object is now defined inside the SDK. (no oratio dependencies)
- NLG object is now defined inside the SDK. (no oratio dependencies)


##  [0.1.0]
###  Added
- Init repository with request and response object
- The server object handle the marshalling and unmarshalling of objects
