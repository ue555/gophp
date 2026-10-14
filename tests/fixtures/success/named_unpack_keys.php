<?php
function capture(...$values) {
    foreach ($values as $key => $value) echo $key, '=', $value, ';';
    echo "\n";
}
capture(...[-4 => 'minus', 8 => 'eight', 'name' => 'text']);
capture(...[], ...[3 => 'three']);
