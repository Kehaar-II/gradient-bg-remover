#!/bin/python3

from PIL import Image
import sys
import os

def clamp(min, x, max):
    '''int -> int'''
    if (x > max):
        return max
    if (x < min):
        return min
    return x

def remove_bg(img, white_level, black_level):
    '''Image, int, int -> Image'''
    map = [[-1 for x in range(img.size[0])] for y in range(img.size[1])]
    alpha = Image.new("L", img.size, 255)

    w, h = img.size
    for x in range(w):
        for y in range(h):
            pixel = img.getpixel((x, y))
            avg = int((pixel[0] + pixel[1] + pixel[2]) / 3)
            if(avg >= black_level):
                alpha.putpixel((x, y), 0)
            if(avg <= white_level):
                alpha.putpixel((x, y), 255)
            alpha.putpixel((x, y), clamp(0, white_level - avg, 255))

    new = Image.new("RGB", img.size, (0, 0, 0))
    new.putalpha(alpha)
    return new

def printHelp():
    file = open("./help", "r")
    print(file.read().replace("%progName%", sys.argv[0]), end = '')
    file.close

def parse_arguments():
    '''void -> Input_Arguments, bool'''
    argv = sys.argv
    argc = len(argv)

    if argc <= 1:
        printHelp()
        return True, None, 0, 0
    if argv[1] == "-h" or argv[1] == "--help":
        printHelp()
        return False, None, 0, 0
    filepath = argv[1]

    if argc == 2:
        return False, filepath, 255, 0

    if argc == 3:
        printHelp()
        return True, None, 0, 0

    if argc == 4:
        return False, filepath, int(argv[2]), int(argv[3])

def main():
    '''void -> int'''
    has_error, filepath, high, low = parse_arguments()
    if (filepath == None):
        return 1 if has_error else 0

    img = Image.open(filepath)
    img = remove_bg(img, high, low)
    img.save(os.path.splitext(filepath)[0] + ".clean" + os.path.splitext(filepath)[1])
    return 0

exit(main())
