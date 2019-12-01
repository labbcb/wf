version 1.0

workflow SayHello {

    input {
        String name = "World"
    }

    call Hello {
        input:
            name = name
    }

    output {
        File msg = Hello.msg
    }
}

task Hello {

  input {
    String name
  }

  command {
    sleep 10
    echo Hello ~{name}! > out.txt
  }

  output {
    File msg = "out.txt"
  }
}
