<?php declare(strict_types=1);
class IntegratedPromoted {
    public function __construct(public readonly int $first, public readonly int $second) {
        echo 'body|';
    }
}
function integratedMark(string $label, $value) { echo $label, '|'; return $value; }
try {
    new IntegratedPromoted(integratedMark('first', 1), integratedMark('second', 'bad'));
} catch (TypeError $e) { echo 'type|'; }
$ok = new IntegratedPromoted(2, 3);
echo $ok->first, '|', $ok->second;
