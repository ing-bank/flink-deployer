import org.apache.flink.api.java.utils.ParameterTool
import org.apache.flink.streaming.api.scala._

object WordCountStateful {

  /** Main program method */
  def main(args: Array[String]) : Unit = {
    // get the execution environment
    val env: StreamExecutionEnvironment = StreamExecutionEnvironment.getExecutionEnvironment

    var intervalMs = 1000
    try {
      val params = ParameterTool.fromArgs(args)
      intervalMs = if (params.has("interval")) params.getInt("interval") else 1000
    } catch {
      case e: Exception => {
        System.err.println("No interval specified. Please run 'WordCountStateful " +
          "--intervalMs <intervalMs>'")
        System.err.println(e)
        return
      } 
    }
    
    // get input data by connecting to the socket
    val text: DataStream[String] = env.addSource(new InfiniteSource(intervalMs = intervalMs))

    // parse the data, group it, window it, and aggregate the counts 
    val windowCounts = text
          .flatMap { w => w.split("\\s") }
          .map { w => WordWithCount(w, 1) }
          .keyBy("word")
          .sum("count")

    // print the results with a single thread, rather than in parallel
    windowCounts.print().setParallelism(1)

    env.execute("Windowed WordCount")
  }

  /** Data type for words with count */
  case class WordWithCount(word: String, count: Long)
}
