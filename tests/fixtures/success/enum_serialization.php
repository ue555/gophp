<?php
enum SerialState: string { case Ready = 'ready'; case Done = 'done'; }
$case = SerialState::Ready;
$encoded = serialize($case);
$decoded = unserialize($encoded);
echo $encoded, '|', $decoded === $case ? 'same' : 'different', '|';
$collection = unserialize(serialize([$case, SerialState::Done]));
echo $collection[0] === $case ? 'same' : 'different', '|', $collection[1]->value;
