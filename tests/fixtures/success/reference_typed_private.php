<?php
class TypedPrivate {
    private int $value = 1;
    public function link() { $alias =& $this->value; $alias = '8'; }
    public function read() { return $this->value; }
}
$box = new TypedPrivate;
$box->link();
echo $box->read(), '|';
try { $outside =& $box->value; } catch (Error $e) { echo 'blocked'; }
