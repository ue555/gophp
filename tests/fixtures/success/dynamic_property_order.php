<?php
class OrderBox { public int $value = 1; }
function receiver() { global $box; echo 'receiver|'; return $box; }
function member() { echo 'member|'; return 'value'; }
function rhs() { echo 'rhs|'; return 5; }
$box = new OrderBox;
echo (receiver()->{member()} = rhs()), '|', $box->value, "\n";
