import selenium.webdriver
import time, requests,json,random,os, faker,names
from random import uniform
from collections import OrderedDict
from faker import Faker
from faker.config import AVAILABLE_LOCALES
def main():
    if not os.path.exists(os.path.join(os.getcwd(),"downloadPics")):
        os.mkdir(os.path.join(os.getcwd(),"downloadPics"))
    numUsers = int(input("Number of users: ") )
    numPics = int(input("Number of pictures per user: ") )

    DRIVER_PATH = 'chromedriver.exe'
    wd = selenium.webdriver.Chrome(executable_path=DRIVER_PATH)
    fake = Faker(AVAILABLE_LOCALES)
    for i in range(numUsers):
        userURL = "http://localhost:8080/users"
        name = names.get_full_name()
        args = {}
        args["username"] = name
        args["password"] = "password"
        args["nickname"] = name

        #Register the new user
        jsonParams = json.dumps(args)
        requests.post(userURL, data = jsonParams)
        header = {"content-type":"application/json"}
        #Login as new user
        args.pop("nickname")
        jsonParams = json.dumps(args)
        ret = requests.post(userURL + "/login", data=jsonParams, headers=header)
        #Get the token
        token = (json.loads(ret.content))["token"]

        #Generate image urls
        locale = AVAILABLE_LOCALES[random.randint(0,len(AVAILABLE_LOCALES)-1)]
        fileList = []
        country = fake[locale].country()
        city = fake[locale].city()
        state = fake[locale].city()
        
        query = {"scenery","buildings","nature","landmarks"}

        urls = fetch_image_urls(country + " " + random.choice(tuple(query)),numPics,wd)
        for url in urls:
            filename = "downloadPics/"+str(random.randint(0,99999999999999999)) + '.jpg'
            fileList.append(filename)
            with open(filename, 'wb') as handle:
                response = requests.get(url, stream=True)
                if not response.ok:
                    print(response)
                for block in response.iter_content(1024):
                    if not block:
                        break
                    handle.write(block)

        for fileName in fileList:
            header = {"authorization": token }
            files = {'upload': open(fileName, 'rb')}
            userURL = "http://localhost:8080"
            ret = requests.post(userURL + "/pictures/upload", files=files, headers=header)
            imageID = (json.loads(ret.content))["img"]
            #longi, lat = uniform(-180,180), uniform(-90, 90)
            picData = {
                'title': country + " " + city + " " + state,
                'Img': imageID,
                'lat': float(fake[locale].latitude()),
                'lng': float(fake[locale].longitude()),
                'location' : {
                    'country':country,
                    'state':state,
                    'city':city
                    },
                'Iso': 400,
                'FocalLength':2,
                'Apertur':2.5,
                'ShutterSpeed':0.24,
                'Timestamp':2147483647,
                'Orientation': 24,
                'Elevation': 123.43
                }
            jsonParams = json.dumps(picData)
            headers={'Authorization': token}
            ret = requests.post(userURL + "/pictures/", data=jsonParams, headers=headers)





#code from here
#https://medium.com/@wwwanandsuresh/web-scraping-images-from-google-9084545808a2
def fetch_image_urls(query:str, max_links_to_fetch:int, wd:selenium.webdriver, sleep_between_interactions:int=1):
    def scroll_to_end(wd):
        wd.execute_script("window.scrollTo(0, document.body.scrollHeight);")
        time.sleep(sleep_between_interactions)    
    
    # build the google query
    search_url = "https://www.google.com/search?safe=off&site=&tbm=isch&source=hp&q={q}&oq={q}&gs_l=img"

    # load the page
    wd.get(search_url.format(q=query))

    image_urls = set()
    image_count = 0
    results_start = 0
    while image_count < max_links_to_fetch:
        scroll_to_end(wd)

        # get all image thumbnail results
        thumbnail_results = wd.find_elements_by_css_selector("img.Q4LuWd")
        number_results = len(thumbnail_results)
        
        print(f"Found: {number_results} search results. Extracting links from {results_start}:{number_results}")
        
        for img in thumbnail_results[results_start:number_results]:
            # try to click every thumbnail such that we can get the real image behind it
            try:
                img.click()
                time.sleep(sleep_between_interactions)
            except Exception:
                continue

            # extract image urls    
            actual_images = wd.find_elements_by_css_selector('img.n3VNCb')
            for actual_image in actual_images:
                if actual_image.get_attribute('src') and 'http' in actual_image.get_attribute('src'):
                    image_urls.add(actual_image.get_attribute('src'))

            image_count = len(image_urls)

            if len(image_urls) >= max_links_to_fetch:
                print(f"Found: {len(image_urls)} image links, done!")
                break
        else:
            print("Found:", len(image_urls), "image links, looking for more ...")
            time.sleep(30)
            return
            load_more_button = wd.find_element_by_css_selector(".mye4qd")
            if load_more_button:
                wd.execute_script("document.querySelector('.mye4qd').click();")

        # move the result startpoint further down
        results_start = len(thumbnail_results)

    return image_urls




main()

