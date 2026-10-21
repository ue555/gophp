<?php
class TypedClosureBox {
    public int $value = 1;

    public function reader() {
        return function () { return $this->value; };
    }
}

$box = new TypedClosureBox;
$alias =& $box->value;
$reader = $box->reader();
unset($box);

try {
    $alias = [];
    echo 'accepted|';
} catch (TypeError $e) {
    echo 'rejected|';
}
echo gettype($reader());
