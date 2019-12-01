version 1.0

import "hello.wdl" as SayHello

workflow SayHelloBatch {
    input {
        Array[String] names
    }

    scatter (name in names) {
        call SayHello.SayHello {
            input:
                name = name
        }
    }

    output {
        Array[File] msgs = SayHello.msg
    }
}