import tensorflow as tf
import base64
import simple_pb2
import random


def parse_Simple(string_tensor):
    fields = []
    types = []
    out_names = []

    fields.append("floatval1")
    types.append(tf.float32)
    out_names.append("floatval1")

    fields.append("floatval2")
    types.append(tf.float32)
    out_names.append("floatval2")

    fields.append("floatval3")
    types.append(tf.float32)
    out_names.append("floatval3")


    descriptor_source = base64.b64decode(b'CusCCgxzaW1wbGUucHJvdG8iYgoGU2ltcGxlEhwKCWZsb2F0dmFsMRgBIAEoAlIJZmxvYXR2YWwxEhwKCWZsb2F0dmFsMhgCIAEoAlIJZmxvYXR2YWwyEhwKCWZsb2F0dmFsMxgDIAEoAlIJZmxvYXR2YWwzQghaBnNpbXBsZUrkAQoGEgQAAAgBCggKAQwSAwAAEAoICgEIEgMCABsKCQoCCAsSAwIAGwoKCgIEABIEBAAIAQoKCgMEAAESAwQIDgoLCgQEAAIAEgMFBBgKDAoFBAACAAUSAwUECQoMCgUEAAIAARIDBQoTCgwKBQQAAgADEgMFFhcKCwoEBAACARIDBgQYCgwKBQQAAgEFEgMGBAkKDAoFBAACAQESAwYKEwoMCgUEAAIBAxIDBhYXCgsKBAQAAgISAwcEGAoMCgUEAAICBRIDBwQJCgwKBQQAAgIBEgMHChMKDAoFBAACAgMSAwcWF2IGcHJvdG8z')

    _, outputs = tf.io.decode_proto(string_tensor, "Simple", fields, types, b'bytes://' + descriptor_source)


    return {out_names[i]: outputs[i] for i in range(len(fields))}



def main():
    
    
    Simple = simple_pb2.Simple()

    Simple.floatval1 = random.random()

    Simple.floatval2 = random.random()

    Simple.floatval3 = random.random()

    val = Simple.SerializeToString()
    print(parse_Simple(val))



if __name__ == "__main__":
    main()