$n = 2..10000
$pi = 3
$tick = 0
for ($i = 0; $i -lt $n.count - 1; $i += 2) {
    if ($tick % 2) {
        $pi -= 4 / ($n[$i] * $n[$i+1] * $n[$i+2])
    } else {
        $pi += 4 / ($n[$i] * $n[$i+1] * $n[$i+2])
    }
    $tick++
}
$pi
