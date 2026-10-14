<?php
class CallableBox {
    public function add($a, $b) { return $a + $b; }
    public static function twice($n) { return $n * 2; }
}
$box = new CallableBox;
$instance = [$box, 'add'];
$static = ['CallableBox', 'twice'];
echo $instance(2, 3), '|', $static(4), "\n";
