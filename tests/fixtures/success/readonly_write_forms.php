<?php
class ReadonlyWriteForms { public function __construct(public readonly int $value, public readonly array $items) {} }
$o = new ReadonlyWriteForms(2, [1]);
try { $o->value += 3; } catch (Error $e) { echo 'compound|'; }
try { ++$o->value; } catch (Error $e) { echo 'increment|'; }
try { $o->value--; } catch (Error $e) { echo 'decrement|'; }
try { unset($o->items[0]); } catch (Error $e) { echo 'array-unset|'; }
$o->value ??= 7;
echo $o->value, '|', $o->items[0], "\n";
