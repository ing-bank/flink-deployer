import org.apache.flink.streaming.api.functions.source.SourceFunction
import org.apache.flink.streaming.api.functions.source.SourceFunction.SourceContext

class InfiniteSource(intervalMs: Long) extends SourceFunction[String] with Serializable {
  private var isRunning: Boolean = true

  def run(ctx: SourceContext[String]) = {
    while (isRunning) {
      ctx.markAsTemporarilyIdle()
      Thread.sleep(intervalMs)
      ctx.collect("this")
    }
  }

  def cancel() = isRunning = false
}
