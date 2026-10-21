<?php
$first = new Exception('first');
$second = new Error('second', 4, $first);
$third = new Exception(message: 'third', previous: $second);
echo ($first->getPrevious() === null ? 'null' : 'wrong'), '|';
echo ($second->getPrevious() === $first ? 'same' : 'wrong'), '|';
echo ($third->getPrevious() === $second ? 'same' : 'wrong'), '|';
echo $third->getPrevious()->getPrevious()->getMessage(), '|', $second->getCode(), "\n";
