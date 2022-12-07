package waBinary

//import byteBufferPool "ws/client/plugin/byte_buffer_pool"
//
//// AcquireBinaryEncoder returns an empty byte buffer from the pool.
////
//// Got byte buffer may be returned to the pool via Put call.
//// This reduces the number of memory allocations required for byte buffer
//// management.
//func AcquireBinaryEncoder() BinaryEncoder { return newEncoder(byteBufferPool.AcquireBuffer()) }
//
//// ReleaseBinaryEncoder returns byte buffer to the pool.
////
//// ByteBuffer.B mustn't be touched after returning it to the pool.
//// Otherwise data races will occur.
//func ReleaseBinaryEncoder(encoder BinaryEncoder) { byteBufferPool.ReleaseBuffer(encoder.B) }
