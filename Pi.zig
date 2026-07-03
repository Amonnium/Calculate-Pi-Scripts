const std = @import("std");

pub fn main() void {
    std.debug.print("Value of PI: {}\n", .{std.math.pi});
}
