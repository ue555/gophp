<?php
$plain = static function ($number) { return $number + 1; };
$copy = Closure::bind($plain, null, null);
echo $copy(2), '|';
class BindBox { public int $value = 7; }
$get = function () { return $this->value; };
$bound = Closure::bind($get, new BindBox, BindBox::class);
echo $bound();
