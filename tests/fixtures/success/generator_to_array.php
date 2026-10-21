<?php
function sample() {
    yield 'a' => 2;
    yield 5 => 3;
}
$array = iterator_to_array(sample());
echo $array['a'], '|', $array[5];
