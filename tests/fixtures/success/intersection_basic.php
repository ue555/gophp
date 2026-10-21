<?php
interface IntersectionA {}
interface IntersectionB {}
class IntersectionBoth implements IntersectionA, IntersectionB {}
function intersectionTake(IntersectionA&IntersectionB $value): IntersectionB&IntersectionA { return $value; }
$both = new IntersectionBoth;
echo gettype(intersectionTake($both)), '|', intersectionTake($both) === $both ? 'same' : 'copy', "\n";
