<?php
class HttpFactory {
    public static string $type = HttpProduct::class;
    public static array $types = ['default' => HttpProduct::class];
}
class HttpProduct {
    public function __construct(public int $amount) {}
}
$a = new HttpFactory::$type(3);
$b = new HttpFactory::$types['default'](5);
$diff = array_udiff_assoc(['same' => 1, 'other' => 2], ['same' => 1], static fn ($x, $y) => $x <=> $y);
echo $a->amount, '|', $b->amount, '|', implode(',', array_keys($diff));
