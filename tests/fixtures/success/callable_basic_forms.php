<?php
function callableSum(int $a, int $b = 1): int { return $a + $b; }
class CallableBox {
    public static function twice(int $n): int { return $n * 2; }
    public function plus(int $n): int { return $n + 3; }
}
$box = new CallableBox;
$function = 'callableSum';
$static = 'CallableBox::twice';
$pair = [$box, 'plus'];
$classPair = ['CallableBox', 'twice'];
echo $function(2, 4), '|', $static(5), '|', $pair(6), '|', $classPair(7), '|';
echo (int)is_callable($function), (int)is_callable($static), (int)is_callable($pair), (int)is_callable('missingCallable');
