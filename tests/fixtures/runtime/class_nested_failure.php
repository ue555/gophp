<?php
class Nested {
    function inner() { return 1 / 0; }
    function outer() { return $this->inner(); }
}
$o = new Nested;
$o->outer();
