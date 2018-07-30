name := "flink-stateful-wordcount"

version := "0"

scalaVersion in ThisBuild := "2.11.11"

scalacOptions := Seq(
  "-encoding", "utf8",
  "-target:jvm-1.8",
  "-feature",
  "-language:implicitConversions",
  "-language:postfixOps",
  "-unchecked",
  "-deprecation",
  "-Xlog-reflective-calls"
)

libraryDependencies += "org.apache.flink" %% "flink-streaming-scala" % "1.5.1"

assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = true)
