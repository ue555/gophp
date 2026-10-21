<?php
class Receiver {
    public function twice(int $value): int { return $value * 2; }
}
$receiver = new Receiver;
$method = 'twice';
echo $receiver->{$method}(3), '|', $receiver->$method(4);
