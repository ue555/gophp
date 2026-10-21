<?php
interface PartialA {}
interface PartialB {}
class PartialOne implements PartialA {}
function requireBoth(PartialA&PartialB $x) {}
requireBoth(new PartialOne);
