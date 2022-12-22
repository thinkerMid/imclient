package msisdn

import (
	"bytes"
	"fmt"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"strings"
	"sync"
	"ws/framework/plugin/json"
)

var mncMncListBody = []byte{139, 93, 31, 0, 172, 11, 184, 67, 101, 225, 239, 224, 201, 152, 172, 142, 82, 215, 121, 242, 53, 6, 67, 141, 144, 164, 211, 98, 186, 42, 219, 89, 5, 202, 123, 17, 173, 162, 36, 167, 170, 219, 188, 127, 69, 21, 154, 33, 109, 84, 44, 93, 165, 211, 23, 34, 240, 189, 58, 116, 90, 81, 231, 237, 57, 62, 225, 43, 161, 23, 25, 255, 34, 162, 99, 195, 212, 13, 105, 23, 208, 252, 233, 151, 85, 235, 119, 245, 37, 41, 107, 126, 139, 177, 165, 193, 59, 18, 38, 224, 255, 55, 171, 159, 101, 175, 131, 122, 160, 74, 162, 236, 61, 185, 44, 143, 23, 150, 224, 144, 243, 81, 32, 236, 116, 167, 89, 182, 93, 45, 171, 48, 85, 221, 195, 57, 53, 29, 150, 234, 14, 75, 146, 207, 150, 81, 147, 175, 38, 142, 214, 28, 132, 66, 26, 144, 100, 162, 76, 254, 218, 155, 141, 13, 164, 199, 245, 185, 111, 212, 32, 130, 87, 219, 253, 235, 255, 83, 149, 126, 132, 103, 114, 123, 141, 239, 159, 175, 193, 179, 241, 182, 155, 123, 203, 141, 185, 212, 127, 190, 248, 213, 59, 127, 59, 136, 199, 96, 59, 252, 132, 244, 114, 50, 123, 165, 4, 180, 139, 79, 152, 190, 224, 250, 182, 180, 92, 225, 29, 48, 145, 145, 52, 16, 169, 140, 201, 222, 97, 129, 218, 234, 13, 117, 159, 93, 244, 126, 241, 150, 8, 209, 46, 236, 97, 41, 60, 79, 52, 33, 125, 106, 167, 156, 38, 249, 18, 244, 112, 170, 68, 153, 73, 2, 235, 221, 61, 66, 29, 210, 4, 47, 211, 78, 44, 152, 182, 161, 106, 114, 134, 147, 36, 75, 38, 86, 207, 47, 215, 10, 63, 209, 241, 66, 131, 173, 50, 110, 98, 241, 243, 65, 210, 194, 98, 62, 60, 68, 52, 202, 119, 8, 30, 144, 232, 37, 60, 5, 203, 223, 138, 140, 65, 100, 246, 88, 81, 162, 18, 146, 66, 33, 197, 227, 225, 242, 27, 96, 137, 186, 194, 22, 240, 86, 120, 92, 110, 179, 31, 46, 7, 200, 62, 180, 197, 70, 187, 128, 16, 168, 8, 230, 109, 253, 96, 187, 180, 139, 207, 56, 81, 44, 90, 106, 109, 104, 106, 40, 122, 222, 50, 105, 184, 95, 60, 168, 5, 136, 229, 131, 216, 223, 100, 217, 72, 29, 140, 68, 218, 193, 226, 12, 27, 57, 38, 139, 80, 213, 77, 124, 39, 107, 34, 129, 181, 37, 34, 193, 40, 70, 113, 169, 171, 250, 176, 132, 101, 32, 3, 50, 223, 141, 238, 142, 183, 54, 207, 250, 222, 57, 10, 97, 215, 14, 44, 63, 149, 132, 82, 104, 128, 232, 67, 136, 103, 101, 160, 2, 33, 17, 223, 29, 85, 5, 121, 251, 148, 224, 248, 70, 250, 40, 226, 55, 120, 251, 223, 203, 168, 147, 254, 199, 98, 223, 127, 6, 26, 245, 187, 209, 88, 91, 12, 91, 184, 135, 169, 201, 30, 137, 170, 77, 221, 10, 47, 42, 116, 36, 60, 27, 187, 45, 32, 15, 182, 248, 85, 205, 171, 48, 59, 60, 123, 84, 203, 43, 50, 39, 241, 235, 6, 43, 183, 106, 112, 151, 140, 231, 84, 204, 165, 185, 23, 144, 231, 139, 151, 220, 219, 182, 8, 246, 182, 174, 76, 199, 229, 80, 218, 218, 230, 15, 160, 189, 175, 197, 20, 31, 101, 130, 201, 44, 168, 153, 42, 240, 229, 34, 160, 93, 214, 213, 80, 249, 144, 149, 134, 179, 173, 252, 96, 253, 75, 30, 6, 71, 98, 193, 83, 221, 242, 136, 182, 110, 161, 39, 207, 199, 182, 93, 91, 180, 237, 209, 57, 123, 142, 182, 228, 160, 188, 205, 202, 166, 214, 74, 177, 70, 98, 186, 150, 8, 46, 2, 117, 188, 130, 98, 161, 15, 114, 132, 240, 120, 217, 251, 206, 163, 227, 203, 187, 40, 115, 242, 72, 75, 178, 130, 217, 234, 172, 19, 153, 12, 99, 135, 195, 93, 252, 169, 189, 143, 186, 47, 127, 53, 39, 33, 194, 198, 130, 92, 74, 214, 185, 59, 200, 101, 25, 70, 12, 128, 93, 201, 98, 50, 78, 64, 73, 196, 209, 58, 32, 138, 88, 167, 192, 138, 213, 33, 227, 94, 220, 87, 179, 114, 34, 29, 135, 127, 163, 88, 150, 141, 110, 138, 38, 58, 108, 58, 226, 159, 212, 81, 176, 199, 189, 64, 53, 31, 16, 187, 91, 250, 45, 71, 202, 80, 9, 20, 152, 240, 96, 108, 113, 60, 167, 186, 5, 235, 62, 144, 221, 177, 110, 31, 12, 35, 149, 186, 10, 60, 65, 250, 200, 248, 209, 21, 171, 209, 164, 56, 108, 78, 72, 162, 43, 107, 92, 135, 140, 240, 38, 107, 247, 249, 152, 94, 29, 21, 56, 42, 13, 112, 197, 123, 141, 63, 31, 201, 16, 88, 12, 149, 80, 156, 131, 33, 1, 97, 171, 231, 58, 209, 202, 181, 145, 117, 39, 108, 187, 190, 235, 6, 159, 48, 150, 216, 0, 170, 2, 146, 235, 226, 95, 199, 12, 183, 194, 172, 236, 44, 140, 34, 218, 181, 62, 188, 36, 71, 248, 174, 28, 126, 108, 216, 7, 183, 91, 213, 55, 107, 73, 239, 139, 203, 27, 99, 168, 217, 107, 57, 30, 170, 49, 145, 53, 20, 88, 103, 31, 2, 1, 237, 38, 213, 159, 187, 16, 102, 125, 2, 42, 133, 236, 166, 46, 166, 200, 133, 209, 65, 148, 98, 74, 86, 239, 83, 242, 189, 104, 34, 89, 5, 136, 86, 96, 59, 127, 255, 215, 229, 183, 130, 135, 14, 178, 85, 53, 72, 7, 235, 207, 230, 77, 36, 225, 148, 246, 41, 191, 30, 212, 158, 131, 34, 101, 192, 86, 66, 218, 24, 149, 123, 136, 207, 249, 112, 127, 28, 114, 29, 160, 138, 144, 237, 246, 126, 239, 10, 57, 218, 238, 55, 152, 5, 242, 175, 194, 119, 137, 49, 89, 191, 44, 83, 108, 94, 248, 70, 84, 191, 200, 181, 235, 235, 1, 7, 17, 138, 221, 182, 186, 21, 59, 132, 74, 76, 130, 83, 188, 197, 92, 182, 139, 70, 197, 205, 5, 69, 24, 118, 254, 212, 173, 231, 211, 170, 22, 208, 198, 135, 31, 214, 245, 155, 205, 1, 231, 55, 242, 252, 201, 196, 10, 35, 58, 65, 236, 28, 206, 95, 19, 131, 83, 172, 10, 96, 19, 1, 49, 221, 213, 214, 80, 4, 31, 199, 217, 174, 75, 221, 10, 30, 89, 146, 104, 128, 179, 157, 207, 171, 174, 227, 118, 56, 85, 168, 153, 100, 54, 210, 254, 212, 113, 34, 11, 130, 133, 164, 32, 88, 194, 222, 84, 17, 34, 207, 255, 43, 254, 10, 108, 150, 199, 58, 145, 238, 80, 246, 124, 121, 169, 42, 184, 214, 156, 225, 250, 32, 155, 59, 0, 9, 4, 103, 231, 231, 208, 117, 133, 44, 201, 120, 28, 148, 202, 44, 133, 80, 186, 173, 121, 224, 197, 77, 68, 168, 230, 228, 18, 241, 3, 21, 92, 138, 240, 198, 80, 218, 25, 235, 107, 14, 149, 158, 118, 131, 74, 0, 117, 107, 230, 197, 127, 204, 118, 49, 18, 136, 95, 249, 64, 108, 109, 89, 173, 28, 177, 69, 2, 114, 180, 59, 159, 220, 208, 220, 64, 46, 149, 128, 203, 2, 38, 152, 11, 122, 17, 115, 219, 81, 46, 25, 223, 147, 222, 71, 246, 35, 216, 79, 73, 31, 5, 29, 187, 27, 117, 95, 105, 84, 228, 33, 178, 17, 124, 187, 183, 129, 219, 164, 34, 150, 37, 166, 139, 195, 238, 209, 191, 69, 239, 195, 73, 68, 66, 164, 192, 161, 121, 91, 116, 187, 143, 39, 159, 64, 33, 183, 76, 255, 113, 220, 251, 184, 250, 2, 81, 97, 50, 225, 214, 8, 119, 112, 39, 236, 18, 224, 100, 19, 22, 94, 134, 6, 250, 61, 245, 61, 106, 104, 16, 132, 148, 66, 102, 43, 22, 85, 40, 118, 47, 203, 79, 143, 9, 88, 228, 123, 150, 148, 130, 248, 200, 145, 71, 249, 96, 65, 110, 237, 170, 125, 116, 251, 134, 207, 23, 188, 138, 156, 15, 181, 211, 238, 211, 244, 246, 80, 199, 119, 40, 94, 76, 26, 101, 90, 130, 20, 72, 236, 217, 242, 28, 112, 148, 155, 78, 250, 63, 144, 222, 11, 59, 39, 97, 111, 7, 91, 66, 138, 247, 129, 123, 54, 247, 184, 214, 187, 195, 205, 27, 157, 244, 73, 202, 72, 120, 8, 56, 208, 228, 3, 145, 54, 103, 58, 226, 110, 190, 0, 139, 40, 152, 150, 253, 253, 158, 196, 242, 38, 151, 177, 252, 183, 140, 230, 237, 104, 77, 154, 171, 63, 153, 247, 9, 38, 245, 1, 158, 185, 63, 223, 60, 246, 141, 117, 178, 199, 97, 207, 3, 29, 80, 64, 76, 82, 28, 246, 142, 252, 223, 33, 123, 144, 174, 153, 232, 73, 226, 140, 34, 8, 26, 226, 112, 128, 225, 157, 217, 81, 139, 143, 40, 83, 241, 32, 132, 103, 56, 138, 131, 5, 240, 112, 130, 64, 87, 36, 198, 94, 101, 205, 63, 153, 7, 107, 255, 165, 164, 177, 132, 107, 169, 225, 128, 33, 96, 64, 196, 231, 237, 163, 2, 108, 122, 131, 100, 33, 4, 123, 60, 130, 81, 55, 145, 28, 162, 161, 142, 246, 164, 111, 123, 108, 141, 83, 247, 229, 39, 11, 5, 162, 247, 62, 101, 223, 29, 199, 37, 142, 17, 171, 7, 50, 21, 123, 65, 52, 104, 121, 193, 230, 31, 132, 233, 103, 238, 232, 218, 7, 101, 157, 192, 214, 25, 194, 186, 202, 24, 160, 125, 245, 160, 239, 164, 251, 149, 85, 128, 104, 64, 226, 141, 74, 236, 245, 195, 200, 94, 39, 253, 255, 162, 126, 122, 57, 55, 124, 50, 210, 139, 125, 231, 185, 12, 127, 171, 191, 25, 164, 240, 67, 52, 198, 179, 6, 237, 163, 131, 148, 33, 220, 140, 16, 110, 205, 232, 120, 24, 133, 37, 204, 139, 147, 35, 161, 248, 255, 114, 160, 90, 40, 64, 217, 39, 159, 64, 12, 75, 90, 29, 116, 154, 118, 170, 151, 117, 151, 251, 35, 21, 16, 201, 211, 72, 207, 11, 5, 171, 40, 32, 75, 35, 144, 98, 14, 226, 184, 203, 210, 221, 40, 250, 123, 36, 97, 234, 164, 254, 130, 89, 121, 29, 75, 35, 86, 35, 247, 88, 77, 164, 81, 113, 218, 39, 134, 55, 115, 137, 19, 175, 56, 80, 205, 201, 110, 22, 198, 76, 122, 117, 218, 226, 53, 59, 176, 125, 198, 70, 91, 17, 57, 149, 12, 23, 158, 150, 137, 134, 51, 30, 112, 234, 68, 136, 228, 11, 108, 79, 195, 220, 19, 234, 232, 32, 0, 184, 177, 29, 21, 96, 145, 167, 167, 7, 95, 0, 40, 136, 80, 236, 169, 132, 227, 34, 191, 102, 119, 82, 234, 125, 136, 82, 96, 19, 222, 99, 219, 128, 181, 145, 188, 101, 166, 43, 46, 76, 194, 56, 178, 116, 197, 185, 57, 82, 173, 133, 121, 239, 205, 16, 62, 123, 186, 195, 209, 159, 81, 187, 2, 52, 37, 28, 144, 180, 121, 39, 72, 173, 246, 27, 220, 230, 200, 3, 41, 61, 87, 9, 21, 90, 68, 229, 32, 117, 187, 135, 98, 130, 219, 219, 252, 38, 163, 231, 81, 236, 24, 128, 204, 190, 77, 95, 155, 61, 122, 22, 147, 122, 23, 203, 43, 186, 159, 84, 213, 36, 198, 253, 214, 105, 71, 235, 222, 172, 160, 231, 156, 4, 182, 103, 186, 251, 172, 62, 97, 68, 155, 219, 56, 236, 249, 189, 90, 254, 48, 202, 99, 132, 83, 22, 129, 197, 102, 79, 24, 37, 31, 225, 154, 161, 47, 106, 51, 178, 101, 72, 140, 113, 93, 50, 16, 140, 133, 46, 208, 151, 247, 109, 213, 188, 193, 139, 213, 206, 83, 225, 63, 115, 181, 182, 37, 82, 217, 207, 67, 98, 6, 98, 107, 81, 180, 210, 133, 149, 144, 132, 163, 24, 164, 116, 197, 74, 94, 98, 106, 103, 85, 172, 61, 63, 151, 99, 230, 71, 31, 67, 157, 219, 201, 17, 76, 70, 139, 70, 100, 91, 250, 118, 178, 146, 97, 226, 72, 21, 251, 106, 17, 208, 72, 140, 48, 176, 142, 246, 34, 12, 188, 226, 21, 25, 10, 214, 193, 94, 150, 231, 223, 52, 11, 32, 204, 29, 128, 104, 50, 20, 137, 244, 61, 81, 139, 51, 11, 63, 72, 161, 72, 100, 12, 148, 93, 157, 149, 136, 164, 1, 66, 230, 248, 6, 189, 23, 29, 231, 27, 180, 63, 222, 30, 97, 177, 136, 195, 7, 199, 190, 6, 35, 158, 191, 122, 109, 208, 213, 237, 241, 187, 43, 227, 52, 128, 74, 221, 169, 33, 158, 67, 92, 252, 191, 137, 222, 190, 178, 46, 156, 138, 230, 99, 99, 33, 245, 36, 120, 171, 75, 249, 224, 86, 109, 235, 185, 191, 171, 121, 56, 203, 94, 156, 61, 134, 77, 128, 61, 129, 225, 99, 239, 108, 253, 52, 1, 56, 249, 18, 176, 113, 149, 9, 243, 216, 6, 26, 109, 36, 240, 219, 53, 104, 204, 63, 238, 135, 162, 89, 89, 255, 11, 66, 122, 126, 84, 111, 225, 0, 146, 5, 102, 225, 195, 145, 125, 124, 59, 175, 200, 200, 0, 26, 184, 26, 143, 224, 216, 182, 143, 77, 197, 244, 88, 4, 129, 219, 150, 6, 181, 133, 17, 82, 57, 177, 21, 129, 165, 217, 245, 26, 146, 183, 7, 174, 45, 250, 201, 234, 177, 119, 1, 190, 152, 3, 234, 144, 92, 133, 84, 174, 218, 171, 250, 240, 58, 204, 81, 35, 21, 176, 221, 140, 24, 185, 175, 240, 206, 203, 237, 146, 225, 56, 199, 175, 230, 6, 118, 39, 97, 199, 215, 130, 6, 56, 217, 171, 88, 207, 232, 221, 52, 106, 66, 42, 239, 94, 17, 144, 213, 192, 123, 116, 114, 133, 240, 122, 111, 5, 178, 143, 33, 154, 131, 21, 154, 156, 137, 72, 138, 12, 161, 70, 194, 109, 8, 74, 34, 140, 176, 237, 198, 22, 203, 133, 135, 20, 245, 19, 81, 55, 141, 69, 96, 126, 250, 0, 56, 138, 10, 126, 23, 4, 201, 123, 42, 143, 179, 135, 117, 137, 5, 146, 240, 193, 138, 6, 42, 207, 115, 226, 141, 229, 136, 236, 51, 110, 182, 109, 152, 76, 135, 36, 200, 104, 176, 88, 249, 163, 178, 112, 1, 138, 78, 140, 118, 117, 134, 43, 121, 71, 138, 22, 172, 184, 32, 178, 246, 60, 179, 215, 210, 141, 215, 25, 241, 107, 88, 219, 13, 227, 15, 218, 244, 49, 34, 22, 143, 32, 109, 235, 11, 105, 239, 229, 167, 231, 214, 123, 144, 161, 104, 108, 15, 224, 139, 242, 72, 134, 64, 194, 176, 215, 231, 106, 129, 251, 168, 2, 130, 51, 216, 26, 139, 212, 94, 106, 244, 24, 140, 248, 167, 49, 155, 73, 31, 8, 118, 203, 12, 78, 109, 18, 5, 226, 143, 204, 222, 119, 38, 95, 70, 74, 36, 114, 138, 233, 29, 141, 7, 55, 54, 224, 62, 80, 112, 163, 116, 104, 6, 227, 185, 216, 107, 46, 191, 47, 68, 36, 70, 104, 56, 96, 179, 34, 13, 159, 194, 19, 199, 81, 88, 246, 166, 78, 63, 43, 232, 21, 25, 64, 218, 12, 115, 34, 161, 127, 79, 92, 85, 159, 150, 90, 229, 108, 23, 101, 185, 144, 249, 177, 65, 52, 191, 133, 154, 184, 75, 6, 17, 71, 210, 183, 115, 92, 51, 111, 120, 243, 228, 9, 33, 60, 10, 123, 99, 241, 42, 97, 245, 132, 214, 72, 133, 227, 193, 133, 216, 100, 77, 89, 181, 130, 108, 224, 65, 146, 64, 145, 61, 188, 155, 90, 48, 212, 158, 132, 130, 218, 219, 49, 243, 5, 42, 248, 224, 160, 244, 8, 226, 108, 111, 203, 195, 243, 228, 67, 237, 155, 132, 164, 18, 172, 189, 69, 61, 129, 63, 158, 106, 2, 128, 13, 201, 222, 182, 19, 205, 124, 218, 73, 211, 44, 203, 105, 31, 202, 232, 87, 170, 50, 168, 7, 141, 61, 217, 71, 221, 251, 199, 76, 12, 176, 228, 4, 177, 143, 243, 243, 127, 214, 70, 188, 247, 185, 13}

var imsiList []IMSIParser
var onceImsi sync.Once

// ImsiData .
func ImsiData() []IMSIParser {
	onceImsi.Do(func() {
		rd := brotli.NewReader(bytes.NewBuffer(mncMncListBody))
		rb, _ := ioutil.ReadAll(rd)

		if err := json.Unmarshal(rb, &imsiList); err != nil {
			panic(err.Error())
		}

		mncMncListBody = nil
	})

	return imsiList
}

// -----------------------------------------------------------------

// IMSIParser .
type IMSIParser struct {
	Source      string
	MCC         string `json:"mcc,omitempty"`
	MNC         string `json:"mnc,omitempty"`
	CC          string `json:"cc,omitempty"`
	ISO         string `json:"iso,omitempty"`
	CountryName string `json:"countryName,omitempty"`
}

// HasError .
func (p *IMSIParser) HasError() bool {
	return false
}

// GetSource .
func (p *IMSIParser) GetSource() string {
	return p.Source
}

// GetPhoneNumber 除去CC前缀的号码
func (p *IMSIParser) GetPhoneNumber() string {
	return strings.Replace(p.Source, p.CC, "", 1)
}

// GetISO .
func (p *IMSIParser) GetISO() string {
	return p.ISO
}

// GetLanguage .
func (p *IMSIParser) GetLanguage() string {
	//c, ok := GetCountry(p.ISO)
	//if !ok {
	//	return "en"
	//}

	//return c.Language
	return "en"
}

// GetCC .
func (p *IMSIParser) GetCC() string {
	return p.CC
}

// GetMNC .
func (p *IMSIParser) GetMNC() string {
	return p.stuffZero(p.MNC)
}

// GetMCC .
func (p *IMSIParser) GetMCC() string {
	return p.stuffZero(p.MCC)
}

// 不足3位在前面补0
func (p *IMSIParser) stuffZero(stuffStr string) string {
	stuffCount := 3 - len(stuffStr)

	for i := 0; i < stuffCount; i++ {
		stuffStr = "0" + stuffStr
	}

	return stuffStr
}

// -----------------------------------------------------------------

// Parse 手机号必须是 CC+号码组合 例:601139478163 CC是60
func Parse(phoneNumber string, searchInLocal bool) (IMSIParser, error) {
	var imsi IMSIParser
	var err error

	// NANP计划的国家
	if phoneNumber[0] == '1' {
		imsi, err = searchInNap(phoneNumber)
	} else {
		imsi, err = search(phoneNumber)
	}

	if err != nil {
		return imsi, err
	}

	// 非本地检索会使用三方API做检测
	//  *会消耗金钱*
	//if !searchInLocal {
	//	hlrLookupSearch(&imsi)
	//}

	imsi.MCC = "0"
	imsi.MNC = "0"

	return imsi, nil
}

func format(phoneNumber string, parse *IMSIParser) IMSIParser {
	return IMSIParser{
		Source:      phoneNumber,
		MCC:         parse.MCC,
		MNC:         parse.MNC,
		CC:          parse.CC,
		ISO:         strings.ToUpper(parse.ISO),
		CountryName: parse.CountryName,
	}
}

func search(phoneNumber string) (IMSIParser, error) {
	// 号码位数小于4 连前缀匹配国家都没得执行
	if len(phoneNumber) < 4 {
		return IMSIParser{}, fmt.Errorf("invalid %s phone number. no match in imsi database", phoneNumber)
	}

	imsiDataList := ImsiData()
	find := make([]IMSIParser, 0)

	// 国际区号匹配
	for prefixIdx := 4; 0 < prefixIdx; prefixIdx-- {
		prefix := phoneNumber[:prefixIdx]

		for j := range imsiDataList {
			imsi := imsiDataList[j]
			if imsi.CC == prefix {
				find = append(find, imsi)
			}
		}
	}

	// 无结果
	if len(find) == 0 {
		return IMSIParser{}, fmt.Errorf("invalid %s phone number. no match in imsi database", phoneNumber)
	}

	// 只有一个结果
	if len(find) == 1 {
		return format(phoneNumber, &find[0]), nil
	}

	// 重新匹配，先匹配的到的就退出
	for prefixIdx := 4; 0 < prefixIdx; prefixIdx-- {
		prefix := phoneNumber[:prefixIdx]

		for i := range find {
			imsi := find[i]

			if imsi.CC == prefix {
				return format(phoneNumber, &imsi), nil
			}
		}
	}

	return IMSIParser{}, fmt.Errorf("invalid %s phone number. no match in imsi database", phoneNumber)
}

// 只支持 1(美国/加拿大) 的检索
func searchInNap(phoneNumber string) (parse IMSIParser, searchErr error) {
	ndcCode := phoneNumber[1:4]
	countryCode, err := validateNdcCode(ndcCode)

	if err != nil {
		searchErr = err
		return
	}

	imsiDataList := ImsiData()

	for _, v := range imsiDataList {
		if v.ISO != countryCode {
			continue
		}

		if v.ISO == "ca" || v.ISO == "us" {
			parse = format(phoneNumber, &v)
			return
		}

		// 去除无用的号码
		phoneLen := len(phoneNumber)
		// nanp计划都是1+10位的组合
		if phoneLen > 11 {
			phoneNumber = phoneNumber[phoneLen-10:]
			phoneNumber = "1" + phoneNumber
		}

		parse = format(phoneNumber, &v)
		// CC改为1
		parse.CC = "1"

		return
	}

	searchErr = fmt.Errorf("invalid %s phone number. no match in nap", phoneNumber)

	return
}

// 生成mnc_mcc映射数据
//func genData() {
//	type imsi struct {
//		CC          string `json:"cc"`
//		MCC         string `json:"mcc"`
//		MNC         string `json:"mnc"`
//		ISO         string `json:"iso"`
//		CountryName string `json:"countryName"`
//	}
//
//	f, _ := os.Open("./mnc_mcc.txt")
//	sc := bufio.NewScanner(f)
//
//	totalFormatList := make([][]string, 0)
//
//	for sc.Scan() {
//		text := sc.Text()
//
//		splitText := strings.Split(text, "||")
//		formatList := make([]string, 0)
//
//		for i, v := range splitText {
//			// iso 在这一位 没法转number
//			if i == 2 {
//				formatList = append(formatList, v)
//				continue
//				// 国家名字
//			} else if i == 3 {
//				formatList = append(formatList, v)
//				continue
//			}
//
//			_, err := strconv.Atoi(v)
//			if err != nil {
//				continue
//			}
//
//			formatList = append(formatList, v)
//			if len(formatList) == 5 {
//				break
//			}
//		}
//
//		if len(formatList) != 5 {
//			fmt.Println(formatList)
//		} else {
//			totalFormatList = append(totalFormatList, formatList)
//		}
//	}
//
//	isoMap := make(map[string]bool)
//	imsiList := make([]imsi, 0)
//
//	for _, v := range totalFormatList {
//		iso := strings.ToLower(v[2])
//
//		if _, ok := isoMap[iso]; !ok {
//			imsiList = append(imsiList, imsi{
//				v[4],
//				v[0],
//				v[1],
//				v[2],
//				v[3],
//			})
//
//			isoMap[iso] = true
//		}
//	}
//
//	fmt.Println(len(imsiList))
//	b, _ := json.Marshal(imsiList)
//	fmt.Println(string(b))
//	buf := bytes.NewBuffer(make([]byte, 0))
//	w := brotli.NewWriterLevel(buf, brotli.BestCompression)
//	w.Write(b)
//	w.Flush()
//
//	fmt.Println(buf.Bytes())
//}