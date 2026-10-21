<?php
class TypedClosureArrayBox {
    public int $value = 1;

    public function reader() {
        return function () { return $this->value; };
    }
}

$array = [1];
$box = new TypedClosureArrayBox;
$box->value =& $array[0];
$reader = $box->reader();
unset($box);
$copy = $array;
$copy[0] = 4;
echo $array[0], '|', $copy[0], '|', $reader();
