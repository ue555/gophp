<?php
class C {
    public $p = null; private $secret = 4; protected $hidden = 0;
    function inspect() { echo isset($this->secret) + 0, empty($this->hidden) + 0; }
}
$c = new C;
echo isset($c->p) + 0, isset($c->secret) + 0, empty($c->secret) + 0, empty($c->hidden) + 0;
echo isset($c->unknown) + 0, empty($c->unknown) + 0, ':'; $c->inspect(); echo "\n";
$x = 4; echo isset($x->p) + 0, empty($x->p) + 0, "\n";
