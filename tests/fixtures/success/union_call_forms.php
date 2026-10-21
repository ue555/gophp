<?php
function unionCollect(int|string ...$items): array { return $items; }
class UnionBox {
    public static function staticValue(string|int $value): string|int { return $value; }
    public function instanceValue(int|string $value): int|string { return $value; }
}
$closure = function (int|string $value): string|int { return $value; };
$arrow = fn(float|int $value): int|float => $value;
$box = new UnionBox;
$values = unionCollect(second: 'b', first: 1);
echo count($values), '|', gettype($values['first']), '|', gettype($values['second']), '|';
echo gettype(UnionBox::staticValue(value: 'x')), '|', gettype($box->instanceValue(2)), '|';
echo gettype($closure('y')), '|', gettype($arrow(3)), "\n";
