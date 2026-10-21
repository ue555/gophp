<?php
interface A {}
interface B {}
interface C {}
echo "before|";
function redundant((A&B)|(A&B&C) $value) {}
echo "after";
