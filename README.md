# slim4go
FitNesse Slim server for Go

As any [Slim server](http://fitnesse.org/FitNesse.UserGuide.WritingAcceptanceTests.SliM.SlimProtocol), slim4go waits for Slim requests to arrive, parses and executes them, and sends the responses back.

Dependency injection and interfaces are used to keep things as isolated as possible and with that testable.

See cmd/slim4godemo for an example of how to use.

Package Structure:

![UML Diagram showing packages](http://www.plantuml.com/plantuml/png/ZPHDRiCW48NtdC8NY5TTLxaAnPE8YXzh85Mgo7TlZ4YXpxQcMKNptZTctjYSKzQSRzuf5RIdD6j3W_7Jy533yzTgoLd_TeqJ-LYrlxhNDWoFfIYBMlfsTDT-TfGsFTTc5tlFDrv5eEe3rtfNjI4J1-rgBn0ksfHE0uXwdeavygg1PE8JlEUjK0-s5Mpu95C1J8X2jlbxNtFnkY_C70sb5FbGpj54jwycuY_Q8xCEa-R9sG_MN8vK7Ay0nowmq-czrTiO0DIaq5q60sk929mLbuqrUDdO9f2yxPooiL-8R6yhaBsm4iYivGvKzmhyZ_Xzs_49_MX7Y4nW-2AoFQ-4x0-Fr2jvOTE2kVKN21n2zaDE0C07AgShGNWq6S845gMCdyRkhX_BlRuIhrjyx6_jOtijAbN_b29y7g310b74xqsfCuNfvjqF)
