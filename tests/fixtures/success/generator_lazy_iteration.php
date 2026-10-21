<?php
function numbers() {
    echo 'start|';
    yield 'first' => 2;
    echo 'middle|';
    yield 3;
    echo 'end|';
    return 7;
}
$generator = numbers();
echo 'created|';
foreach ($generator as $key => $value) {
    echo $key, ':', $value, '|';
}
echo $generator->getReturn();
