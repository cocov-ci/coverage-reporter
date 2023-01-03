package com.company.project.batch.parsing

import com.fasterxml.jackson.databind.ObjectMapper
import java.io.BufferedWriter
import java.io.Closeable

class ParsedBatchWriter(private val mapper: ObjectMapper, private val bufferedWriter: BufferedWriter) : Closeable {
  fun write(value: BatchRow) {
    val stringValue = mapper.writeValueAsString(value)
    bufferedWriter.write(stringValue)
    bufferedWriter.newLine()
  }

  override fun close() {
    bufferedWriter.close()
  }
}
