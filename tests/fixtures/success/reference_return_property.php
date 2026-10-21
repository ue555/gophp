<?php
class ReturnBox {
    public int $value = 1;
    public function &valueRef(): int {
        return $this->value;
    }
}
$box = new ReturnBox;
$alias =& $box->valueRef();
$alias = '7';
echo $box->value, '|', gettype($alias), '|';
$copy = $box->valueRef();
$copy = 9;
echo $box->value;
