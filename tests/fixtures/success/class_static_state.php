<?php
class StateBox {
    public static $values = [1];
    public static function add($n) { self::$values[0] += $n; return self::$values[0]; }
}
echo StateBox::add(2), '|', StateBox::add(4), '|';
echo isset(StateBox::$values[0]), '|', empty(StateBox::$values[0]), '|';
unset(StateBox::$values[0]);
echo isset(StateBox::$values[0]), '|';
StateBox::$values[0] ??= 8;
echo StateBox::$values[0], "\n";
