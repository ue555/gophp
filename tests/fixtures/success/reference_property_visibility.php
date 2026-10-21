<?php
class RefVisible {
    protected $value = 1;
    public function change() { $alias =& $this->value; $alias = 4; }
    public function get() { return $this->value; }
}
$box = new RefVisible;
$box->change();
echo $box->get(), '|';
try { $alias =& $box->value; } catch (Error $e) { echo 'blocked'; }
