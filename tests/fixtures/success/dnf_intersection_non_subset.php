<?php
interface NonSubsetA {}
interface NonSubsetB {}
interface NonSubsetC {}
class AB implements NonSubsetA, NonSubsetB {}
class AC implements NonSubsetA, NonSubsetC {}
function distinct((NonSubsetA&NonSubsetB)|(NonSubsetA&NonSubsetC) $value): string { return gettype($value); }
echo distinct(new AB), '|', distinct(new AC), "\n";
