for times in  {1..3}
do
        for ti in  {1..1000}
        do
            echo "statsd.counter.${ti}:1|c" | nc -w 1 -u 127.0.0.1 8125
        done
        sleep 60
done