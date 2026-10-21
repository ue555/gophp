<?php
interface ImplicitA {}
interface ImplicitB {}
function implicitBoth(ImplicitA&ImplicitB $value = null) { echo gettype($value), "\n"; }
implicitBoth();
